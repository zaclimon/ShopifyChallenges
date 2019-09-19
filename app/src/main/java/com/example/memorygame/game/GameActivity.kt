package com.example.memorygame.game

import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.ImageView
import androidx.appcompat.app.AppCompatActivity
import androidx.constraintlayout.widget.ConstraintLayout
import com.example.memorygame.R
import com.example.memorygame.data.Card
import com.example.memorygame.data.Product
import com.example.memorygame.data.ProductList
import com.google.android.material.bottomsheet.BottomSheetBehavior
import com.squareup.moshi.Moshi
import java.util.*

class GameActivity : AppCompatActivity() {

    private var pairCount: Int = 0
    private var flippedCards: MutableList<Card> = mutableListOf()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_game)

        val bottomSheet = findViewById<ConstraintLayout>(R.id.score_results_bottom_sheet)
        val behavior = BottomSheetBehavior.from(bottomSheet)
        behavior.bottomSheetCallback = object : BottomSheetBehavior.BottomSheetCallback() {
            override fun onStateChanged(bottomSheet: View, newState: Int) {
            }
            override fun onSlide(bottomSheet: View, slideOffset: Float) {
                behavior.state = BottomSheetBehavior.STATE_COLLAPSED
            }
        }
        val products = getGameProducts()
        val views = getGameCardViews()
        val cards = createCards(products, views)
        initListeners(cards)
    }

    private fun getGameProducts(): List<Product> {
        val jsonString = application.assets.open("products.json").bufferedReader().use {
            it.readLine()
        }

        val moshi = Moshi.Builder().build()
        val jsonAdapter = moshi.adapter<ProductList>(ProductList::class.java)
        val products = (jsonAdapter.fromJson(jsonString) as ProductList).products
        return products.take(10)
    }

    private fun getGameCardViews(): List<ImageView> {
        val panel = findViewById<ConstraintLayout>(R.id.game_panel)
        val gameCardsCount = panel.childCount
        val gameCardViews = mutableListOf<ImageView>()
        for (i in 0 until gameCardsCount) {
            gameCardViews.add(panel.getChildAt(i) as ImageView)
        }

        for (view in gameCardViews) {
            flippedCards
            view.setOnClickListener {
                Log.d("GameActivity", view.id.toString())
            }
        }
        return gameCardViews
    }

    private fun createCards(products: List<Product>, views: List<ImageView>): List<Card> {
        val cardsList = mutableListOf<Card>()
        val viewsLinkedList = LinkedList<ImageView>(views)

        /*
          There should definitely be a better way of assigning two ImageViews for the same Product
          though this is the only way I have found. Another solution could be to add the products
          based on it's modulo position so we could avoid "replicating" the card instantiation.
        */
        products.forEach {
            cardsList.addAll(listOf(Card(it, viewsLinkedList.pop()), Card(it, viewsLinkedList.pop())))
        }
        return cardsList
    }

    private fun initListeners(cards: List<Card>) {
        for (card in cards) {
            card.view.setOnClickListener {
                pairCount++
                flippedCards.add(card)
                verifyCards()
            }
        }
    }

    private fun verifyCards() {
        if (flippedCards.size == 2) {
            val isSame = flippedCards.all { card -> card.product.id == card.product.id }
            flippedCards.clear()
        }
    }
}
