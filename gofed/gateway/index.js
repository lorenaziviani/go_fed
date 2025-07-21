import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { buildSubgraphSchema } from '@apollo/subgraph';
import { gql } from 'graphql-tag';

// Configuração das variáveis de ambiente
const USERS_SERVICE_URL = process.env.USERS_SERVICE_URL || 'http://localhost:8081/query';
const PRODUCTS_SERVICE_URL = process.env.PRODUCTS_SERVICE_URL || 'http://localhost:8082/query';
const GATEWAY_PORT = process.env.GATEWAY_PORT || 4000;

// Schema federado com @key
const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.0",
          import: ["@key", "@extends", "@external", "@requires", "@provides"])

  type User @key(fields: "id") {
    id: ID!
    name: String!
    email: String!
  }

  type Product @key(fields: "id") {
    id: ID!
    name: String!
    description: String!
    price: Float!
    category: String!
    owner: User!
  }

  type SemaphoreStats {
    max: Int!
    current: Int!
    available: Int!
    usage: Int!
  }

  type Query {
    users: [User!]!
    user(id: ID!): User
    usersByIds(ids: [ID!]!): [User!]!
    products: [Product!]!
    product(id: ID!): Product
    productsByIds(ids: [ID!]!): [Product!]!
    productsByCategory(category: String!): [Product!]!
    productsWithSemaphore(ids: [ID!]!): [Product!]!
    semaphoreStats: SemaphoreStats!
  }
`;

// Resolvers que fazem proxy para os serviços
const resolvers = {
  Query: {
    users: async () => {
      const response = await fetch(USERS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: '{ users { id name email } }' }),
      });
      const data = await response.json();
      return data.data.users;
    },
    user: async (_, { id }) => {
      const response = await fetch(USERS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: `{ user(id: "${id}") { id name email } }` }),
      });
      const data = await response.json();
      return data.data.user;
    },
    usersByIds: async (_, { ids }) => {
      const response = await fetch(USERS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: `{ usersByIds(ids: [${ids.map(id => `"${id}"`).join(', ')}]) { id name email } }`
        }),
      });
      const data = await response.json();
      return data.data.usersByIds;
    },
    products: async () => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: '{ products { id name description price category owner { id name email } } }' }),
      });
      const data = await response.json();
      return data.data.products;
    },
    product: async (_, { id }) => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: `{ product(id: "${id}") { id name description price category owner { id name email } } }` }),
      });
      const data = await response.json();
      return data.data.product;
    },
    productsByIds: async (_, { ids }) => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: `{ productsByIds(ids: [${ids.map(id => `"${id}"`).join(', ')}]) { id name description price category owner { id name email } } }`
        }),
      });
      const data = await response.json();
      return data.data.productsByIds;
    },
    productsByCategory: async (_, { category }) => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: `{ productsByCategory(category: "${category}") { id name description price category owner { id name email } } }`
        }),
      });
      const data = await response.json();
      return data.data.productsByCategory;
    },
    productsWithSemaphore: async (_, { ids }) => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: `{ productsWithSemaphore(ids: [${ids.map(id => `"${id}"`).join(', ')}]) { id name description price category owner { id name email } } }`
        }),
      });
      const data = await response.json();
      return data.data.productsWithSemaphore;
    },
    semaphoreStats: async () => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: '{ semaphoreStats { max current available usage } }' }),
      });
      const data = await response.json();
      return data.data.semaphoreStats;
    },
  },
  // Resolvers para federation
  User: {
    __resolveReference: async (reference) => {
      const response = await fetch(USERS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: `{ user(id: "${reference.id}") { id name email } }` }),
      });
      const data = await response.json();
      return data.data.user;
    },
  },
  Product: {
    __resolveReference: async (reference) => {
      const response = await fetch(PRODUCTS_SERVICE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: `{ product(id: "${reference.id}") { id name description price category owner { id name email } } }` }),
      });
      const data = await response.json();
      return data.data.product;
    },
  },
};

// Criar o servidor Apollo
const server = new ApolloServer({
  schema: buildSubgraphSchema({ typeDefs, resolvers }),
  introspection: true,
});

// Iniciar o servidor
async function startServer() {
  const { url } = await startStandaloneServer(server, {
    listen: { port: parseInt(GATEWAY_PORT) },
  });

  console.log(`Apollo Federation Gateway ready at ${url}`);
  console.log(`GraphQL Playground available at ${url}`);
  console.log('\nConnected Services:');
  console.log(`   - users: ${USERS_SERVICE_URL}`);
  console.log(`   - products: ${PRODUCTS_SERVICE_URL}`);
  console.log('\nFederation is active! You can now query across services.');
  console.log('Federation keys enabled: User(id), Product(id)');
  console.log('\nNew Features:');
  console.log('   - Concurrent user resolution with WaitGroup + Channels');
  console.log('   - Semaphore-limited product resolution (max 3 concurrent)');
  console.log('   - Backpressure control and performance monitoring');
  console.log('\nExample Queries:');
  console.log('   - usersByIds(ids: ["1", "2", "3", "4", "5"])');
  console.log('   - productsWithSemaphore(ids: ["1", "2", "3", "4", "5"])');
  console.log('   - semaphoreStats');
  console.log('   - productsByCategory(category: "Electronics")');
}

startServer().catch(console.error); 