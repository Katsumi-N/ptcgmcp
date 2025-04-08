package com.example

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class SearchResponse(
    val pokemons: List<PokemonData>?,
    val trainers: List<TrainerData>?,
    val energies: List<EnergyData>?
)

@Serializable
data class PokemonData(
    val id: String,
    val name: String,
    val hp: Int?,
    @SerialName("energy_type") val energyType: String?,
    @SerialName("image_url") val imageUrl: String?,
)

@Serializable
data class TrainerData(
    val id: String,
    val name: String,
    @SerialName("trainer_type") val trainerType: String,
    @SerialName("image_url") val imageUrl: String,
)

@Serializable
data class EnergyData(
    val id: String,
    val name: String,
    @SerialName("image_url") val imageUrl: String,
)

sealed class CardDetailResponse {
    data class Pokemon(val details: PokemonDetailResponse) : CardDetailResponse()
    data class Trainer(val details: TrainerDetailResponse) : CardDetailResponse()
    data class Energy(val details: EnergyDetailResponse) : CardDetailResponse()
}


@Serializable
data class PokemonDetailResponse(
    val result: Boolean,
    val pokemon: PokemonDetail,
)

@Serializable
data class PokemonDetail(
    val id: String,
    val name: String,
    val hp: Int,
    @SerialName("image_url") val imageUrl: String?,
    @SerialName("energy_type") val energyType: String,
    val ability: String?,
    @SerialName("ability_description") val abilityDescription: String?,
    val attacks: List<Attack>?,
    val regulation: String?,
    val expansion: String?,
)

@Serializable
data class Attack(
    val name: String,
    @SerialName("required_energy") val requiredEnergy: String,
    val damage: String,
    val description: String?,
)

@Serializable
data class TrainerDetailResponse(
    val result: Boolean,
    val trainer: TrainerDetail,
)

@Serializable
data class TrainerDetail(
    val id: String,
    val name: String,
    @SerialName("trainer_type") val trainerType: String,
    @SerialName("image_url") val imageUrl: String,
    val description: String?,
    val regulation: String?,
    val expansion: String?,
)

@Serializable
data class EnergyDetailResponse(
    val result: Boolean,
    val energy: EnergyDetail,
)

@Serializable
data class EnergyDetail(
    val id: String,
    val name: String,
    @SerialName("image_url") val imageUrl: String,
    val description: String?,
    val regulation: String?,
    val expansion: String?,
)
