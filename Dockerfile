# Multi-stage build for production
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the app
RUN npm run build

# Production stage - serve static files with Node
FROM node:22-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache curl

# Install a lightweight static file server
RUN npm install -g serve

# Copy built assets from builder
COPY --from=builder /app/dist ./dist

# Expose port 3003
EXPOSE 3003

# Serve the static files
CMD ["serve", "-s", "dist", "-l", "3003", "--no-clipboard"]
