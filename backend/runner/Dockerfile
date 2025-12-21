FROM python:3.11-slim

# Create non-root user for security
RUN useradd -m -s /bin/bash runner

# Set working directory
WORKDIR /sandbox

# Create code directory
RUN mkdir -p /sandbox/code && chown runner:runner /sandbox/code

# Switch to non-root user
USER runner

# Default command
CMD ["python3", "-c", "print('ready')"]
