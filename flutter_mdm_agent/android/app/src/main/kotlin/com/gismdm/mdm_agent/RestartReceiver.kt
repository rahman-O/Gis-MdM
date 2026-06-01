package com.gismdm.mdm_agent

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log

/**
 * Receiver triggered by AlarmManager to restart the background service
 * if it was killed by the system.
 *
 * This is a self-healing mechanism: the AgentForegroundService schedules
 * an alarm in onDestroy(). When the alarm fires, this receiver checks
 * if the service is running and restarts it if needed.
 */
class RestartReceiver : BroadcastReceiver() {

    companion object {
        private const val TAG = "MDM_RestartReceiver"
    }

    override fun onReceive(context: Context, intent: Intent) {
        Log.i(TAG, "Restart alarm fired — ensuring service is running")
        startAgentService(context)
    }

    private fun startAgentService(context: Context) {
        val serviceIntent = Intent(context, AgentForegroundService::class.java)
        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                context.startForegroundService(serviceIntent)
            } else {
                context.startService(serviceIntent)
            }
            Log.i(TAG, "Agent service restarted via AlarmManager")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to restart agent service", e)
        }
    }
}
