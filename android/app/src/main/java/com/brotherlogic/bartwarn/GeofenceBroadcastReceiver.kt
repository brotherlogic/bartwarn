package com.brotherlogic.bartwarn

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent

class GeofenceBroadcastReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent) {
        val geofencingEvent = com.google.android.gms.location.GeofencingEvent.fromIntent(intent)
        if (geofencingEvent == null || geofencingEvent.hasError()) {
            android.util.Log.e("GeofenceReceiver", "Error in geofencing event")
            return
        }

        val triggeringGeofences = geofencingEvent.triggeringGeofences
        if (triggeringGeofences != null) {
            val stationIds = triggeringGeofences.map { it.requestId }
            processGeofences(context, geofencingEvent.geofenceTransition, stationIds, System.currentTimeMillis())
        }
    }

    internal fun processGeofences(
        context: Context,
        transitionType: Int,
        stationIds: List<String>,
        currentTime: Long
    ) {
        if (transitionType != com.google.android.gms.location.Geofence.GEOFENCE_TRANSITION_ENTER) {
            return
        }

        val sharedPreferences = context.getSharedPreferences("geofence_cooldown", Context.MODE_PRIVATE)
        val cooldownMillis = 10 * 60 * 1000L // 10 minutes

        for (stationId in stationIds) {
            val lastTriggerTime = sharedPreferences.getLong(stationId, 0L)

            if (currentTime - lastTriggerTime >= cooldownMillis) {
                // Update timestamp
                sharedPreferences.edit().putLong(stationId, currentTime).apply()

                // Enqueue worker
                val constraints = androidx.work.Constraints.Builder()
                    .setRequiredNetworkType(androidx.work.NetworkType.CONNECTED)
                    .build()

                val workRequest = androidx.work.OneTimeWorkRequestBuilder<LocationPingWorker>()
                    .setInputData(androidx.work.workDataOf("station_id" to stationId))
                    .setConstraints(constraints)
                    .build()

                androidx.work.WorkManager.getInstance(context).enqueue(workRequest)
            } else {
                android.util.Log.d("GeofenceReceiver", "Skipping ping for $stationId due to cooldown")
            }
        }
    }
}
