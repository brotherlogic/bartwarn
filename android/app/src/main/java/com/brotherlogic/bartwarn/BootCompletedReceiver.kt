package com.brotherlogic.bartwarn

import android.app.PendingIntent
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log
import com.google.android.gms.location.LocationServices

class BootCompletedReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action == Intent.ACTION_BOOT_COMPLETED) {
            val geofencingClient = LocationServices.getGeofencingClient(context)
            val geofenceManager = GeofenceManager(geofencingClient)

            val geofenceIntent = Intent().setClassName(context.packageName, "com.brotherlogic.bartwarn.GeofenceBroadcastReceiver")
            
            val flags = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
                PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_MUTABLE
            } else {
                PendingIntent.FLAG_UPDATE_CURRENT
            }

            val pendingIntent = PendingIntent.getBroadcast(
                context,
                0,
                geofenceIntent,
                flags
            )

            try {
                geofenceManager.addGeofences(
                    pendingIntent,
                    onSuccess = { Log.d("BootCompletedReceiver", "Geofences re-registered on boot") },
                    onFailure = { e -> Log.e("BootCompletedReceiver", "Failed to re-register geofences on boot", e) }
                )
            } catch (e: SecurityException) {
                Log.e("BootCompletedReceiver", "Missing location permission", e)
            }
        }
    }
}
