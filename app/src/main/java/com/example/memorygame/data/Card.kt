package com.example.memorygame.data

import android.widget.ImageView

/**
 * Holder class representing a Card. Note that we could reduce coupling with the Android framework
 * by using using an ID for example, but for the sake of simplicity and reduced findViewById() calls,
 * it has been added here.
 *
 * @author Isaac Pateau
 */
data class Card(
    val product: Product,
    val view: ImageView
)