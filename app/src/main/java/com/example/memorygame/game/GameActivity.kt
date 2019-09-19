package com.example.memorygame.game

import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AppCompatActivity
import androidx.constraintlayout.widget.ConstraintLayout
import com.example.memorygame.R
import com.google.android.material.bottomsheet.BottomSheetBehavior

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
    }
}
