package com.example.memorygame.data

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

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