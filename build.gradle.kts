
plugins {
    alias(libs.plugins.kotlin.jvm)
    alias(libs.plugins.ktor)
    kotlin("plugin.serialization") version "2.1.20"
    application
}

group = "com.example"
version = "0.0.1"

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(17)
    }
}

application {
    mainClass.set("com.example.MainKt")
}

repositories {
    mavenCentral()
}

val exposedVersion: String by project
val slf4jVersion = "2.0.9"

dependencies {
    implementation("org.slf4j:slf4j-nop:$slf4jVersion")
    testImplementation(libs.kotlin.test.junit)

    // kotest
    testImplementation("io.kotest:kotest-runner-junit5:5.8.1")
    testImplementation("io.kotest:kotest-assertions-core:5.8.1")
    testImplementation("io.kotest.extensions:kotest-assertions-arrow:2.0.0")
    testImplementation("io.kotest:kotest-property:5.8.1")

    // Exposed
    implementation("org.jetbrains.exposed:exposed-core:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-dao:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-jdbc:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-jodatime:$exposedVersion")
    // mysql
    implementation("com.mysql:mysql-connector-j:8.3.0")

    // testcontainers
    testImplementation(platform("org.testcontainers:testcontainers-bom:1.20.6"))
    testImplementation("org.testcontainers:junit-jupiter")
    testImplementation("org.testcontainers:mysql")

    // Arrow
    implementation("io.arrow-kt:arrow-core:2.0.1")
    implementation("io.arrow-kt:arrow-fx-coroutines:2.0.1")

    // MCP
    implementation("io.modelcontextprotocol:kotlin-sdk:0.4.0")
}
