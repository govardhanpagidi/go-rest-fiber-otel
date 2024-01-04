FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o main

# Expose the port your API listens on
EXPOSE 8080

ENV MONGODB_URI="mongodb+srv://admin:admin@fxrate.73lrkis.mongodb.net/"

# Command to run your application
CMD ["./main"]