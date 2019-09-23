package com.zaclimon.memorygame.listeners

/**
 * Interface responsible for listening to events happening during the download process of a Product
 * definition data encoded as a JSON file.
 *
 * @author Isaac Pateau
 */

interface JsonDownloaderListener {

    /**
     * Notifies the user when the JSON definition of one or multiple products have been downloaded.
     */
    fun onJsonDownloadSuccess()

    /**
     * Notifies the user when the JSON definition of one or multiple products have failed.
     */
    fun onJsonDownloadFailed()
}