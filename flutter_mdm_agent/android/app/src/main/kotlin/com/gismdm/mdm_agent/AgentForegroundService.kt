package com.gismdm.mdm_agent

import android.app.AlarmManager
import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.os.IBinder
import android.os.SystemClock
import android.util.Log
import androidx.core.app.NotificationCompat

/**
 * Foreground Service — keeps the MDM Agent alive permanently.
 * Shows a persistent notification to prevent Android from killing the process.
 *
 * Note: The flutter_background_service package provides its own service,
 * but we keep this as a fallback and for AlarmManager self-restart logic.
 * The package's service handles the headless Flutter engine.
 */
class AgentForegroundService : Service() {

    companion object {
        private const val TAG = "MDM_ForegroundService"
        private const val CHANNEL_ID = "mdm_agent_foreground"
        private const val CHANNEL_NAME = "MDM Agent Service"
        private const val NOTIFICATION_ID = 1001
        private const val RESTART_ALARM_REQUEST_CODE = 2001
        private const val RESTART_INTERVAL_MS = 5 * 60 * 1000L // 5 minutes
    }

    override fun onCreate() {
        super.onCreate()
        Log.i(TAG, "Foreground service created")
        createNotificationChannel()
        startForeground(NOTIFICATION_ID, buildNotification())
        scheduleRestartAlarm()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        Log.i(TAG, "Foreground service started")
        return START_STICKY // Restart if killed
    }

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onDestroy() {
        super.onDestroy()
        Log.w(TAG, "Foreground service destroyed — scheduling restart via AlarmManager")
        scheduleRestartAlarm()
    }

    /**
     * Schedule a repeating alarm that will restart the service if it gets killed.
     * Uses AlarmManager.setExactAndAllowWhileIdle() for Doze mode compatibility.
     */
    private fun scheduleRestartAlarm() {
        val alarmManager = getSystemService(Context.ALARM_SERVICE) as AlarmManager
        val restartIntent = Intent(this, RestartReceiver::class.java).apply {
            action = "com.gismdm.mdm_agent.RESTART_SERVICE"
        }
        val pendingIntent = PendingIntent.getBroadcast(
            this,
            RESTART_ALARM_REQUEST_CODE,
            restartIntent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )

        val triggerTime = SystemClock.elapsedRealtime() + RESTART_INTERVAL_MS

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            alarmManager.setExactAndAllowWhileIdle(
                AlarmManager.ELAPSED_REALTIME_WAKEUP,
                triggerTime,
                pendingIntent
            )
        } else {
            alarmManager.setExact(
                AlarmManager.ELAPSED_REALTIME_WAKEUP,
                triggerTime,
                pendingIntent
            )
        }
        Log.i(TAG, "Restart alarm scheduled in ${RESTART_INTERVAL_MS / 1000}s")
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                CHANNEL_NAME,
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "MDM Agent background service"
                setShowBadge(false)
            }
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(channel)
        }
    }

    private fun buildNotification(): Notification {
        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle("MDM Agent")
            .setContentText("Device management active")
            .setSmallIcon(android.R.drawable.ic_lock_idle_lock)
            .setPriority(NotificationCompat.PRIORITY_LOW)
            .setOngoing(true)
            .build()
    }
}
