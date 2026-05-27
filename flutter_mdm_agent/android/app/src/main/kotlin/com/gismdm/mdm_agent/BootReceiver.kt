package com.gismdm.mdm_agent

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log

/**
 * Boot Receiver — starts the Foreground Service when the device boots.
 * Also handles MY_PACKAGE_REPLACED (app updated).
 */
class BootReceiver : BroadcastReceiver() {

    companion object {
        private const val TAG = "MDM_BootReceiver"
    }

    override fun onReceive(context: Context, intent: Intent) {
        val action = intent.action
        Log.i(TAG, "Received broadcast: $action")

        if (action == Intent.ACTION_BOOT_COMPLETED || action == Intent.ACTION_MY_PACKAGE_REPLACED) {
            startAgentService(context)
        }
    }

    private fun startAgentService(context: Context) {
        val serviceIntent = Intent(context, AgentForegroundService::class.java)
        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                context.startForegroundService(serviceIntent)
            } else {
                context.startService(serviceIntent)
            }
            Log.i(TAG, "Agent service started")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to start agent service", e)
        }
    }
}
