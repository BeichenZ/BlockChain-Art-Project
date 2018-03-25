import RPi.GPIO as GPIO
import time

GPIO.setmode(GPIO.BCM)

#Button input pin
GPIO.setup(19,GPIO.IN,pull_up_down=GPIO.PUD_UP)
#Display LED Pin
GPIO.setup(6,GPIO.OUT)

try:
	while True:
		button_state = GPIO.input(19)
		if button_state == False:
			GPIO.output(6,True)
			print('Button Pressed')
			time.sleep(0.2)
		else:
			GPIO.output(6,False)
except:
	GPIO.cleanup()
