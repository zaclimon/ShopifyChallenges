package com.example.memorygame.main

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.FrameLayout
import androidx.constraintlayout.widget.ConstraintLayout
import com.example.memorygame.R
import com.example.memorygame.data.ProductUtils
import com.example.memorygame.game.GameActivity
import com.example.memorygame.listeners.ImageDownloaderListener
import com.example.memorygame.listeners.JsonDownloaderListener
import com.google.android.material.snackbar.Snackbar
import kotlinx.android.synthetic.main.activity_main.*
import java.io.File

class MainActivity : AppCompatActivity(), JsonDownloaderListener, ImageDownloaderListener {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        button_new_game.setOnClickListener {
            val intent = Intent(this, GameActivity()::class.java)
            startActivity(intent)
        }

        // Check if products.json is present so we can download the assets if necessary
        val productFile = File(cacheDir, ProductUtils.PRODUCTS_FILE)

        if (!productFile.isFile) {
            val rootView = findViewById<FrameLayout>(android.R.id.content)
            Snackbar.make(rootView, "Downloading assets for first time use", Snackbar.LENGTH_SHORT).show()
            ProductUtils.saveProductJson(this, this)
        }
    }

    override fun onJsonDownloadSuccess() {
        // Download images when we have the products definition
        ProductUtils.saveProductImages(this, this)
    }

    override fun onJsonDownloadFailed() {
        val rootView = findViewById<FrameLayout>(android.R.id.content)
        Snackbar.make(rootView, "Download failed!", Snackbar.LENGTH_SHORT).show()
    }

    override fun onImagesDownloadSuccess() {
        val rootView = findViewById<FrameLayout>(android.R.id.content)
        Snackbar.make(rootView, "Assets downloaded!", Snackbar.LENGTH_SHORT).show()
    }
}
