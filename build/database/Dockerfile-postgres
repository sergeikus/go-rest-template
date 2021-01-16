FROM postgres:13.1-alpine
ENV POSTGRES_DB=database
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=adminPassword
COPY build/database/scripts/init.sql /docker-entrypoint-initdb.d/
