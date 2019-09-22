package com.zaclimon.memorygame.listeners

/**
 * Interface responsible for listening to events happening during the download process of a Product
 * image.
 *
 * @author Isaac Pateau
 */
interface ImageDownloaderListener {

    /**
     * Notifies the user when the image of one or multiple products have been downloaded.
     */
    fun onImagesDownloadSuccess()
}