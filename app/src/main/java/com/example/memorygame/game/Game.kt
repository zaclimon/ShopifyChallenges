package com.example.memorygame.game

interface Game {

    fun isValidPair(): Boolean
    fun getConnectedPairs(): Int

}