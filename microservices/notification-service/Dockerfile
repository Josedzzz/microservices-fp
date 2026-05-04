# -------- Build stage --------
FROM maven:3.9-eclipse-temurin-17 AS build

WORKDIR /app

COPY pom.xml ./
RUN mvn -B -DskipTests dependency:go-offline

COPY src ./src
RUN mvn -B -DskipTests clean package

# -------- Runtime stage --------
FROM eclipse-temurin:17-jre

WORKDIR /app

RUN groupadd --system spring && useradd --system --gid spring spring

COPY --from=build /app/target/notification-service-*.jar /app/app.jar

USER spring:spring

EXPOSE 8084
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
