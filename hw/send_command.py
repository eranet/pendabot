#!/usr/bin/env python
import asyncio
import random
import json

from nats.aio.client import Client as NATS
import odrive
from odrive.enums import *

async def run(loop):
    nc = NATS()

    await nc.connect(loop=loop)

    async def message_handler(msg): # wait for command and set command to motor
        subject = msg.subject
        reply = msg.reply
        data = json.loads(msg.data.decode())
        print("Received a message on '{subject} {reply}': {data}".format(subject=subject, reply=reply, data=data))
        print(data['Value'])
        # my_drive.axis0.controller.current_setpoint = -data.Value/0.123  # coefficient torque / current

    sid = await nc.subscribe("/pendabot/shoulder_torque_controller/command", cb=message_handler)

    while True: # publish current state
        c_value = random.random()  # my_drive.axis0.motor.current_control.Iq_measured
        data = {
            "Value": c_value,
        }
        await nc.publish('/current', json.dumps(data).encode())
        await asyncio.sleep(0.1)

    # Remove interest in subscription.
    await nc.unsubscribe(sid)

    # Terminate connection to NATS.
    await nc.close()


if __name__ == '__main__':
    # my_drive.reboot()

    # my_drive = odrive.find_any()
    print("ODrive found.")
    # my_drive.axis0.controller.config.control_mode = CTRL_MODE_CURRENT_CONTROL
    # my_drive.axis0.requested_state = AXIS_STATE_CLOSED_LOOP_CONTROL
    # my_drive.axis0.controller.current_setpoint = 0

    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    loop.run_forever()
    loop.stop()
    loop.close()
    # my_drive.reboot()
