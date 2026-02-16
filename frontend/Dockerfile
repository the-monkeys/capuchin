FROM node:24-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install


# Build Stage
FROM node:24-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules .
COPY . .

# Build the application
RUN npm run build

FROM node:24-alpine AS dev

WORKDIR /app

COPY --from=deps /app/node_modules .
COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev"]


FROM nginx:alpine as prod

# Copy the built assets from the builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
