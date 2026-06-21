package com.brotherlogic.bartwarn

import android.app.PendingIntent
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log
import com.google.android.gms.location.LocationServices

class BootCompletedReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action == Intent.ACTION_BOOT_COMPLETED) {
            Log.d("BootCompletedReceiver", "ACTION_BOOT_COMPLETED received. Re-registering geofences.")
            
            val geofencingClient = LocationServices.getGeofencingClient(context)
            val geofenceManager = GeofenceManager(geofencingClient)
            
            val geofenceIntent = Intent("com.brotherlogic.bartwarn.ACTION_GEOFENCE_EVENT")
            geofenceIntent.setPackage(context.packageName)
            
            val pendingIntent = PendingIntent.getBroadcast(
                context,
                0,
                geofenceIntent,
                PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_MUTABLE
            )
            
            geofenceManager.addGeofences(
                pendingIntent = pendingIntent,
                onSuccess = {
                    Log.d("BootCompletedReceiver", "Geofences re-registered successfully.")
                },
                onFailure = { e ->
                    Log.e("BootCompletedReceiver", "Failed to re-register geofences", e)
                }
            )
        }
    }
}
