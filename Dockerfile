# Use the official PostgreSQL image as the base image
FROM postgres:latest

# Set environment variables for PostgreSQL
ENV POSTGRES_DB=your_database_name
ENV POSTGRES_USER=your_username
ENV POSTGRES_PASSWORD=your_password

# Copy the SQL scripts to initialize the database
# COPY init.sql /docker-entrypoint-initdb.d/

# Set the default port for PostgreSQL
EXPOSE 5432