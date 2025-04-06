package com.example

import io.modelcontextprotocol.kotlin.sdk.*
import io.modelcontextprotocol.kotlin.sdk.server.Server
import io.modelcontextprotocol.kotlin.sdk.server.ServerOptions
import io.modelcontextprotocol.kotlin.sdk.server.StdioServerTransport
import kotlinx.coroutines.runBlocking
import kotlinx.io.asSink
import kotlinx.io.asSource
import kotlinx.io.buffered
import kotlinx.coroutines.Job
import kotlinx.serialization.json.*

fun `run mcp server`() {
    val server = configureServer()

    runMcpServerUsingStdio(server)
}

fun configureServer(): Server {
    val server = Server(
        Implementation(
            name = "pokemon trading card game mcp server",
            version = "0.0.1",
        ),
        ServerOptions(
            capabilities = ServerCapabilities(
                prompts = ServerCapabilities.Prompts(listChanged = true),
                resources = ServerCapabilities.Resources(subscribe = true, listChanged = true),
                tools = ServerCapabilities.Tools(listChanged = true),
            )
        )
    )

    server.addTool(
        name = "get_pokemon_card_list",
        description = "ポケモンカードのリストを取得するツール",
        inputSchema = Tool.Input(
            properties = buildJsonObject {
                putJsonObject("name") {
                    put("type", "string")
                    put("description", "ポケモンカードの名前")
                }
            },
            required = listOf("name")
        )
    ) { request ->
        val name = request.arguments["name"]?.jsonPrimitive?.toString()

        CallToolResult(
            content = listOf(TextContent(name))
        )
    }
    server.addResource(
        uri = "file:///Users/naya.katsumi/develop/ptcgmcp/src/main/resources/pokemon_card.csv",
        name = "pokemon_card",
        description = "ポケモンカードのCSVファイル",
        mimeType = "text/csv"
    ) { request ->
        ReadResourceResult(
            contents = listOf(
                TextResourceContents(
                    text = "ポケモンカードの結果",
                    uri = request.uri,
                    mimeType = "text/plain",
                )
            )
        )
    }

    return server
}

fun runMcpServerUsingStdio(server: Server) {
    val transport = StdioServerTransport(
        inputStream = System.`in`.asSource().buffered(),
        outputStream = System.out.asSink().buffered()
    )

    runBlocking {
        server.connect(transport)
        val done = Job()
        server.onClose {
            done.complete()
        }
        done.join()
        println("Server closed")
    }
}
