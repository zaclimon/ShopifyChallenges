package com.example.memorygame.game

import android.os.Bundle
import android.view.View
import android.widget.ImageView
import android.widget.LinearLayout
import androidx.appcompat.app.AppCompatActivity
import androidx.constraintlayout.widget.ConstraintLayout
import com.bumptech.glide.Glide
import com.example.memorygame.R
import com.example.memorygame.data.Card
import com.example.memorygame.data.Product
import com.example.memorygame.data.ProductList
import com.google.android.material.bottomsheet.BottomSheetBehavior
import com.squareup.moshi.Moshi
import kotlinx.android.synthetic.main.bottom_sheet.*
import java.util.*

class GameActivity : AppCompatActivity() {

    private val MAX_PAIRS: Int = 10
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

        reset_button.setOnClickListener { initGame() }
        initGame()
    }

    /**
     * Initializes/Resets everything for a new game
     */
    private fun initGame() {

        // Get all products, ImageViews and generate cards based on them
        val products = getGameProducts()
        val views = getGameCardViews()
        val cards = createCards(products, views)

        // If we were playing previously, reset the pair count
        if (pairCount == MAX_PAIRS) {
            pairCount = 0
        }

        // Return all the flipped cards as their "back" side
        for (view in views) {
            view.setImageResource(R.drawable.ic_shopify)
        }

        // Initialize the listeners when pressing on a card
        initCardListeners(cards)

        // Set the score
        textview_score.text = getString(R.string.score_text, pairCount, MAX_PAIRS)

    }

    /**
     * Retrieves a list of products based on the JSON file given for the challenge
     */
    private fun getGameProducts(): List<Product> {
        val jsonString = application.assets.open("products.json").bufferedReader().use {
            it.readLine()
        }

        val moshi = Moshi.Builder().build()
        val jsonAdapter = moshi.adapter<ProductList>(ProductList::class.java)
        val products = (jsonAdapter.fromJson(jsonString) as ProductList).products
        return products.take(10)
    }

    /**
     * Gets all the views which will be our cards
     */
    private fun getGameCardViews(): List<ImageView> {
        val panel = findViewById<ConstraintLayout>(R.id.game_panel)
        val gameCardsCount = panel.childCount
        val gameCardViews = mutableListOf<ImageView>()
        for (i in 0 until gameCardsCount) {
            val cardLayout = panel.getChildAt(i) as LinearLayout
            gameCardViews.add(cardLayout.getChildAt(0) as ImageView)
        }

        return gameCardViews
    }

    /**
     * Creates the cards based on the products received from Shopify's JSON file.
     */
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

    /**
     * Initializes listeners when touching on the cards.
     */
    private fun initCardListeners(cards: List<Card>) {
        for (card in cards) {
            card.view.setOnClickListener {

                /*
                 Every time we press on a card we "flip" it by doing the following:

                 1. We add the card to the list of currently checked cards
                 2. We replace the image for this card, thus "flipping" it
                 3. We verify to see if the cards are the same
                 */

                flippedCards.add(card)
                Glide.with(this).load(card.product.image.src).into(card.view)
                verifyCards()
            }
        }
    }

    /**
     * Verifies flipped cards to ensure that they are the same. If this is not the case,
     * then they are being flipped to their "back" side.
     */
    private fun verifyCards() {
        // Do nothing unless we have two flipped cards
        if (flippedCards.size == 2) {

            /*
             * When we are sure that we have the same cards, be sure to not consider them anymore by
             * removing their listeners. Also notify the user that we got a new pair! 🎉
             */

            if (flippedCards[0].product.id == flippedCards[1].product.id) {
                pairCount++
                flippedCards.forEach { it.view.setOnClickListener(null) }
                textview_score.text = getString(R.string.score_text, pairCount, MAX_PAIRS)
            } else {
                // Set the card back to the Shopify logo when we don't have a good fit
                flippedCards.forEach { it.view.setImageResource(R.drawable.ic_shopify) }
            }
            flippedCards.clear()
        }
    }
}
