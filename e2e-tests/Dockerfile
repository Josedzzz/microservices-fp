FROM node:20-alpine

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy project files
COPY . .

# Default command to run Cucumber tests
CMD ["./node_modules/.bin/cucumber-js"]
