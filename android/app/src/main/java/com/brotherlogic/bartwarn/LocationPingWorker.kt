package com.brotherlogic.bartwarn

import android.content.Context
import android.util.Log
import androidx.work.CoroutineWorker
import androidx.work.WorkerParameters
import brotherlogic.bartwarn.api.LocationRequest
import brotherlogic.bartwarn.api.LocationServiceGrpcKt
import io.grpc.ManagedChannelBuilder

class LocationPingWorker(appContext: Context, workerParams: WorkerParameters) :
    CoroutineWorker(appContext, workerParams) {

    override suspend fun doWork(): Result {
        val stationId = inputData.getString("station_id") ?: return Result.failure()

        return try {
            // Channel should ideally be injected or managed centrally, but creating it here for the ping
            val channel = ManagedChannelBuilder.forAddress("backend.bartwarn.com", 443)
                .useTransportSecurity()
                .build()

            val stub = LocationServiceGrpcKt.LocationServiceCoroutineStub(channel)
            val request = LocationRequest.newBuilder()
                .setStationId(stationId)
                .build()

            stub.recordLocation(request)
            channel.shutdown()

            Result.success()
        } catch (e: Exception) {
            Log.e("LocationPingWorker", "Error pinging location", e)
            Result.retry()
        }
    }
}
