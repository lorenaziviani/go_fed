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

  type Query {
    users: [User!]!
    user(id: ID!): User
    usersByIds(ids: [ID!]!): [User!]!
    products: [Product!]!
    product(id: ID!): Product
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
  console.log('\n Connected Services:');
  console.log(`   - users: ${USERS_SERVICE_URL}`);
  console.log(`   - products: ${PRODUCTS_SERVICE_URL}`);
  console.log('\n Federation is active! You can now query across services.');
  console.log('Federation keys enabled: User(id), Product(id)');
  console.log('\n Example queries:');
  console.log('  - products { id name owner { id name } }');
  console.log('  - users { id name } products { id name owner { id name } }');
  console.log('  - usersByIds(ids: ["1", "2", "3"]) { id name email }');
}

startServer().catch(console.error); 