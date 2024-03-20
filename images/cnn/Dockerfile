# Use the official PyTorch image as the base image
FROM pytorch/pytorch:latest

# Install any additional dependencies
RUN pip install scikit-learn

# Copy the current directory contents into the container at /app
COPY . /app

# Set the working directory to /app
WORKDIR /app

# Run the Python script when the container launches
CMD ["python", "ml_app.py"]