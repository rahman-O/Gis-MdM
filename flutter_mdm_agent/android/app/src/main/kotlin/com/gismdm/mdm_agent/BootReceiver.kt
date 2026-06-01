package com.gismdm.mdm_agent

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log

/**
 * Boot Receiver — starts the Foreground Service when the device boots.
 * Also handles MY_PACKAGE_REPLACED (app updated).
 *
 * Note: flutter_background_service handles autoStartOnBoot internally,
 * but we keep this receiver as a fallback to ensure the service always starts.
 */
class BootReceiver : BroadcastReceiver() {

    companion object {
        private const val TAG = "MDM_BootReceiver"
    }

    override fun onReceive(context: Context, intent: Intent) {
        val action = intent.action
        Log.i(TAG, "Received broadcast: $action")

        if (action == Intent.ACTION_BOOT_COMPLETED ||
            action == Intent.ACTION_MY_PACKAGE_REPLACED ||
            action == "android.intent.action.QUICKBOOT_POWERON"
        ) {
            startAgentService(context)
        }
    }

    private fun startAgentService(context: Context) {
        // Start our custom foreground service as a fallback
        val serviceIntent = Intent(context, AgentForegroundService::class.java)
        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                context.startForegroundService(serviceIntent)
            } else {
                context.startService(serviceIntent)
            }
            Log.i(TAG, "Agent foreground service started on boot")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to start agent foreground service", e)
        }
    }
}
