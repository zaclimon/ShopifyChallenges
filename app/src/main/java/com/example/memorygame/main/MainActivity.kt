package com.example.memorygame.main

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Button
import com.example.memorygame.R
import com.example.memorygame.game.GameActivity
import kotlinx.android.synthetic.main.activity_main.*

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        button_new_game.setOnClickListener {
            val intent = Intent(this, GameActivity()::class.java)
            startActivity(intent)
        }
    }
}
