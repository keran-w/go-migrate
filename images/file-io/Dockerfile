# Base image
FROM ubuntu:latest

# Update apt and install necessary packages
RUN apt-get update && apt-get install -y g++

# Copy source code into the container
COPY io.cpp /io.cpp

# Copy the host file into the container
COPY /source.txt /source.txt

# Compile the source code
RUN g++ -o /io /io.cpp

# Execute the program
CMD ["/io"]
