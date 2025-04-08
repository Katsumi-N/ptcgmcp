
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
val ktorVersion = "3.1.1"

dependencies {
    implementation("org.slf4j:slf4j-nop:$slf4jVersion")
    testImplementation(libs.kotlin.test.junit)

    // ktor client
    implementation("io.ktor:ktor-client-content-negotiation:$ktorVersion")
    implementation("io.ktor:ktor-serialization-kotlinx-json:$ktorVersion")

    // MCP
    implementation("io.modelcontextprotocol:kotlin-sdk:0.4.0")

}