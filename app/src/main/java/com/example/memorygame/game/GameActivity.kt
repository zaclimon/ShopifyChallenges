package com.example.memorygame.game

import android.os.Bundle
import android.util.Log
import android.view.View
import androidx.appcompat.app.AppCompatActivity
import androidx.constraintlayout.widget.ConstraintLayout
import com.example.memorygame.R
import com.example.memorygame.data.ProductList
import com.google.android.material.bottomsheet.BottomSheetBehavior
import com.squareup.moshi.Moshi

class GameActivity : AppCompatActivity() {

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
        initListeners()
        initProducts()
    }

    private fun initListeners() {
        val panel = findViewById<ConstraintLayout>(R.id.game_panel)
        val gameCardsCount = panel.childCount
        val gameCardViews = mutableListOf<View>()
        for (i in 0 until gameCardsCount) {
            gameCardViews.add(panel.getChildAt(i))
        }

        for (view in gameCardViews) {
            view.setOnClickListener { Log.d("GameActivity", view.id.toString()) }
        }
    }

    private fun initProducts() {
        val jsonString = application.assets.open("products.json").bufferedReader().use {
            it.readLine()
        }

        val moshi = Moshi.Builder().build()
        val jsonAdapter = moshi.adapter<ProductList>(ProductList::class.java)
        val products = (jsonAdapter.fromJson(jsonString) as ProductList).products
        Log.d(localClassName, products.toString())
    }
}
