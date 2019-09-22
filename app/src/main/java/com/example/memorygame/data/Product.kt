package com.example.memorygame.data

import android.content.Context
import com.bumptech.glide.Glide
import com.example.memorygame.listeners.ImageDownloaderListener
import com.example.memorygame.listeners.JsonDownloaderListener
import com.squareup.moshi.JsonClass
import com.squareup.moshi.Moshi
import okhttp3.*
import java.io.File
import java.io.IOException


@JsonClass(generateAdapter = true)
data class Product(
    val id: Long,
    val title: String,
    val vendor: String,
    val image: ProductImage
)

@JsonClass(generateAdapter = true)
data class ProductList(
    val products: List<Product>
)

@JsonClass(generateAdapter = true)
data class ProductImage(
    val id: Long,
    val product_id: Long,
    val width: Int,
    val height: Int,
    val src: String
)

class ProductUtils {

    companion object {

        const val PRODUCTS_URL = "https://shopicruit.myshopify.com/admin/products.json?page=1&access_token=c32313df0d0ef512ca64d5b336a0d7c6"
        const val PRODUCTS_FILE = "products.json"

        /**
         * Retrieves a list of products based on the JSON file given for the challenge
         */
        fun retrieveProducts(context: Context): List<Product> {

            val jsonFile = File(context.cacheDir, PRODUCTS_FILE)
            val jsonString = jsonFile.inputStream().bufferedReader().use {
                it.readLine()
            }

            val moshi = Moshi.Builder().build()
            val jsonAdapter = moshi.adapter<ProductList>(ProductList::class.java)
            return (jsonAdapter.fromJson(jsonString) as ProductList).products
        }

        /**
         * Downloads the products JSON file and saves them in the device's cache folder
         */
        fun saveProductJson(context: Context, listener: JsonDownloaderListener) {

            // Download JSON first
            val client = OkHttpClient.Builder().build()
            val request = Request.Builder()
                .url(PRODUCTS_URL)
                .build()
            client.newCall(request).enqueue(object: Callback {

                // Notify the user if the download has failed
                override fun onFailure(call: Call, e: IOException) {

                    // Let also the developer know what was the cause of the failure
                    e.printStackTrace()
                    listener.onJsonDownloadFailed()
                }

                // Notify the user if the download has been successful
                override fun onResponse(call: Call, response: Response) {

                    // Save the JSON to the cache directory of the application
                    val downloadedFile = File(context.cacheDir, PRODUCTS_FILE)
                    val inputStream = response.body?.byteStream()

                    inputStream?.bufferedReader()?.use {
                        downloadedFile.writeText(it.readText())
                    }

                    inputStream?.close()
                    listener.onJsonDownloadSuccess()
                }
            })
        }

        /**
         * Downloads each products images and then saves them into Glide's cache for further usage.
         */
        fun saveProductImages(context: Context, listener: ImageDownloaderListener) {

            /*
              This is NOT the best way of handling this mainly because we assume that:

              1. The device has reliable internet connection. (Which might be true in a sandbox like
              emulator but not with a real device)
              2. That Shopify's CDN is always available to retrieve the images.

              Also, for the sake of simplicity, a listener for asynchronous based Glide requests has
              not been added since we would need to verify that all the images have been downloaded
              successfully. (Glide doesn't let people easily access it's cache it seems)

              With that said, for a case like this game, I found this to be "relatively" acceptable
              compromises.
             */

            val products = retrieveProducts(context)
            products.forEach {
                Glide.with(context)
                    .downloadOnly()
                    .load(it.image.src)
                    .submit()
            }

            listener.onImagesDownloadSuccess()
        }
    }

}