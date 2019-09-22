package com.zaclimon.memorygame.main

import android.content.Context
import android.content.Intent
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.os.Bundle
import android.widget.FrameLayout
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import com.google.android.material.snackbar.Snackbar
import com.zaclimon.memorygame.BuildConfig
import com.zaclimon.memorygame.R
import com.zaclimon.memorygame.data.ProductUtils
import com.zaclimon.memorygame.game.GameActivity
import com.zaclimon.memorygame.listeners.ImageDownloaderListener
import com.zaclimon.memorygame.listeners.JsonDownloaderListener
import kotlinx.android.synthetic.main.activity_main.*
import java.io.File

/**
 * Activity representing the home screen of the application
 *
 * @author Isaac Pateau
 */
class MainActivity : AppCompatActivity(), JsonDownloaderListener, ImageDownloaderListener {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        button_new_game.setOnClickListener {
            val intent = Intent(this, GameActivity()::class.java)
            startActivity(intent)
        }

        button_about.setOnClickListener {
            AlertDialog.Builder(this)
                .setTitle(R.string.about_text)
                .setMessage(getString(R.string.about_dialog_text, BuildConfig.VERSION_NAME))
                .show()
        }

        // Check if products.json is present so we can download the assets if necessary
        val productFile = File(cacheDir, ProductUtils.PRODUCTS_FILE)

        if (!productFile.isFile) {
            val rootView = findViewById<FrameLayout>(android.R.id.content)
            button_new_game.isEnabled = false

            // Verify if the network is available before letting the user start a new game
            if (isNetworkAvailable()) {
                Snackbar.make(rootView, getString(R.string.download_assets), Snackbar.LENGTH_SHORT).show()
                ProductUtils.saveProductJson(this, this)
            } else {
                Snackbar.make(rootView, getString(R.string.verify_internet_connection), Snackbar.LENGTH_INDEFINITE)
                    .setAction(R.string.dismiss_text) {}
                    .show()
            }
        }
    }

    override fun onJsonDownloadSuccess() {
        // Download images when we have the product definitions
        ProductUtils.saveProductImages(this, this)
    }

    override fun onJsonDownloadFailed() {
        val rootView = findViewById<FrameLayout>(android.R.id.content)
        Snackbar.make(rootView, getString(R.string.download_assets_failure), Snackbar.LENGTH_SHORT).show()
    }

    override fun onImagesDownloadSuccess() {

        /*
          One of the hiccups of running asynchronous tasks without knowing exactly when all images
          will appear is that modifying some UI elements might not be possible because we were
          running on a background thread. "Force" this by explicitly running on the UI thread.
         */
        runOnUiThread {
            val rootView = findViewById<FrameLayout>(android.R.id.content)
            button_new_game.isEnabled = true
            Snackbar.make(rootView, getString(R.string.download_assets_success), Snackbar.LENGTH_SHORT).show()
        }
    }

    /**
     * Verifies if the network is able to access internet so it is possible to download the required
     * assets
     */
    private fun isNetworkAvailable(): Boolean {

        val connectivityManager = getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        val networkCapabilities = connectivityManager.getNetworkCapabilities(connectivityManager.activeNetwork)

        if (networkCapabilities != null) {
            return networkCapabilities.hasCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET)
        }

        return false
    }
}
