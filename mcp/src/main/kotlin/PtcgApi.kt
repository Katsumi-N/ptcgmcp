package com.example

import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json

class PtcgApiClient {
    private val apiBaseUrl = System.getenv("PTCG_API_BASE_URL") ?: "http://localhost:8080"
    private val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json(Json {
                ignoreUnknownKeys = true
                isLenient = true
            })
        }
    }

    suspend fun searchCards(query: String, cardType: String?): SearchResponse {
        val response: SearchResponse = client.get("$apiBaseUrl/v1/cards/search") {
            parameter("q", query)
            parameter("card_type", cardType)
        }.body()

        return response
    }

    suspend fun getCardDetail(cardId: String, cardType: String): CardDetailResponse {
        return when (cardType) {
            "pokemon" -> {
                val response: PokemonDetailResponse = client.get("$apiBaseUrl/v1/cards/detail/pokemon/$cardId").body()
                CardDetailResponse.Pokemon(response)
            }
            "trainer" -> {
                val response: TrainerDetailResponse = client.get("$apiBaseUrl/v1/cards/detail/trainer/$cardId").body()
                CardDetailResponse.Trainer(response)
            }
            "energy" -> {
                val response: EnergyDetailResponse = client.get("$apiBaseUrl/v1/cards/detail/energy/$cardId").body()
                CardDetailResponse.Energy(response)
            }
            else -> throw IllegalArgumentException("Invalid card type: $cardType")
        }
    }

    fun close() {
        client.close()
    }
}
