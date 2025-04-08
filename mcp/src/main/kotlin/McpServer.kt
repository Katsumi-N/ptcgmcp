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
    val apiClient = PtcgApiClient()
    val server = configureServer(apiClient)
    runMcpServerUsingStdio(server, apiClient)
}

fun configureServer(apiClient: PtcgApiClient): Server {
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
        name = "search_pokemon_card",
        description = "ポケモンカードをキーワード検索",
        inputSchema = Tool.Input(
            properties = buildJsonObject {
                putJsonObject("query") {
                    put("type", "string")
                    put("description", "検索するポケモンカードの名前")
                }
                putJsonObject("card_type") {
                    put("type", "string")
                    put("description", "カードの種類を指定 (pokemon | trainer | energy)")
                }
            },
            required = listOf("query")
        )
    ) { request ->
        val query = request.arguments["query"]?.jsonPrimitive?.content
            ?: return@addTool CallToolResult(
                content = listOf(TextContent("query is required"))
            )
        val cardType = request.arguments["card_type"]?.jsonPrimitive?.content

        try {
            val response = apiClient.searchCards(query, cardType)
            val contentList = mutableListOf<TextContent>()

            contentList.add(TextContent("ポケモンカードの検索結果：$query"))

            // ポケモンカードの結果を表示
            if (!response.pokemons.isNullOrEmpty()) {
                contentList.add(TextContent("【ポケモン】"))
                response.pokemons.forEach { pokemon ->
                    contentList.add(TextContent("ID: ${pokemon.id}, 名前: ${pokemon.name}, タイプ: ${pokemon.energyType ?: "情報なし"}, HP: ${pokemon.hp ?: "情報なし"}"))
                }
            }

            // トレーナーカードの結果を表示
            if (!response.trainers.isNullOrEmpty()) {
                contentList.add(TextContent("【トレーナー】"))
                response.trainers.forEach { trainer ->
                    contentList.add(TextContent("ID: ${trainer.id}, 名前: ${trainer.name}, タイプ: ${trainer.trainerType}"))
                }
            }

            // エネルギーカードの結果を表示
            if (!response.energies.isNullOrEmpty()) {
                contentList.add(TextContent("【エネルギー】"))
                response.energies.forEach { energy ->
                    contentList.add(TextContent("ID: ${energy.id}, 名前: ${energy.name}"))
                }
            }

            CallToolResult(content = contentList)
        } catch (e: Exception) {
            CallToolResult(
                content = listOf(TextContent("エラーが発生しました: ${e.message}"))
            )
        }
    }

    server.addTool(
        name = "get_card_detail",
        description = "ポケモンカードの詳細情報を取得",
        inputSchema = Tool.Input(
            properties = buildJsonObject {
                putJsonObject("id") {
                    put("type", "string")
                    put("description", "検索するポケモンカードのID")
                }
                putJsonObject("card_type") {
                    put("type", "string")
                    put("description", "カードの種類を指定 (pokemon | trainer | energy)")
                }
            },
            required = listOf("id", "card_type")
        )
    ) { request ->
        val id = request.arguments["id"]?.jsonPrimitive?.content
            ?: return@addTool CallToolResult(
                content = listOf(TextContent("id is required"))
            )
        val cardType = request.arguments["card_type"]?.jsonPrimitive?.content
            ?: return@addTool CallToolResult(
                content = listOf(TextContent("card_type is required"))
            )

        if (cardType !in listOf("pokemon", "trainer", "energy")) {
            return@addTool CallToolResult(
                content = listOf(TextContent("card_type must be pokemon, trainer or energy. given: $cardType"))
            )
        }

        val response = apiClient.getCardDetail(id, cardType)
        val contentList = mutableListOf<TextContent>()

        try {
            when (response) {
                is CardDetailResponse.Pokemon -> {
                    val pokemon = response.details.pokemon
                    contentList.add(TextContent("ポケモンカードの詳細情報：${pokemon.name}"))
                    contentList.add(TextContent("ID: ${pokemon.id}"))
                    contentList.add(TextContent("HP: ${pokemon.hp}"))
                    contentList.add(TextContent("タイプ: ${pokemon.energyType}"))
                    if (pokemon.ability != null) {
                        contentList.add(TextContent("特性: ${pokemon.ability}"))
                        contentList.add(TextContent("特性の詳細: ${pokemon.abilityDescription}"))
                    }
                    contentList.add(TextContent("ワザ ${pokemon.attacks}"))

                    if (pokemon.regulation != null) {
                        contentList.add(TextContent("レギュレーション: ${pokemon.regulation}"))
                    }
                    if (pokemon.expansion != null) {
                        contentList.add(TextContent("拡張パック: ${pokemon.expansion}"))
                    }
                }
                is CardDetailResponse.Trainer -> {
                    val trainer = response.details.trainer
                    contentList.add(TextContent("トレーナーカードの詳細情報：${trainer.name}"))
                    contentList.add(TextContent("ID: ${trainer.id}"))
                    contentList.add(TextContent("トレーナータイプ: ${trainer.trainerType}"))
                    if (trainer.description != null) {
                        contentList.add(TextContent("効果: ${trainer.description}"))
                    }
                    if (trainer.regulation != null) {
                        contentList.add(TextContent("レギュレーション: ${trainer.regulation}"))
                    }
                    if (trainer.expansion != null) {
                        contentList.add(TextContent("拡張パック: ${trainer.expansion}"))
                    }
                }
                is CardDetailResponse.Energy -> {
                    val energy = response.details.energy
                    contentList.add(TextContent("エネルギーカードの詳細情報：${energy.name}"))
                    contentList.add(TextContent("ID: ${energy.id}"))
                    if (energy.description != null) {
                        contentList.add(TextContent("効果: ${energy.description}"))
                    }
                    if (energy.regulation != null) {
                        contentList.add(TextContent("レギュレーション: ${energy.regulation}"))
                    }
                }
            }

            CallToolResult(content = contentList)
        } catch (e: Exception) {
            return@addTool CallToolResult(
                content = listOf(TextContent("エラーが発生しました: ${e.message}"))
            )
        }
    }

    return server
}

fun runMcpServerUsingStdio(server: Server, apiClient: PtcgApiClient) {
    val transport = StdioServerTransport(
        inputStream = System.`in`.asSource().buffered(),
        outputStream = System.out.asSink().buffered()
    )

    runBlocking {
        server.connect(transport)
        val done = Job()
        server.onClose {
            apiClient.close()
            done.complete()
        }
        done.join()
        println("Server closed")
    }
}
