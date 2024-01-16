# Find the PID of the process
pid=$(pgrep -f "./go_echo")

# Check if the PID is not empty
if [ -n "$pid" ]; then
    echo "Found PID: $pid"

    # Terminate the process using the found PID
    kill "$pid" && go build && nohup ./go_echo &

    echo "Process terminated. and Restarted"
else
    echo "Process not found."
fi