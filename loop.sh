#!/bin/bash
# Usage: ./loop.sh [plan] [max_iterations]
# Examples:
#   ./loop.sh              # Build mode, unlimited iterations
#   ./loop.sh 20           # Build mode, max 20 iterations
#   ./loop.sh plan         # Plan mode, unlimited iterations
#   ./loop.sh plan 5       # Plan mode, max 5 iterations

# Define colors for terminal output
CYAN='\033[36m'
YELLOW='\033[33m'
GREEN='\033[32m'
RED='\033[31m'
GRAY='\033[90m'
NC='\033[0m' # No Color

# Parse arguments
if [ "$1" = "plan" ]; then
  # Plan mode
  MODE="plan"
  PROMPT_FILE="PROMPT_plan.md"
  MAX_ITERATIONS=${2:-0}
elif [[ "$1" =~ ^[0-9]+$ ]]; then
  # Build mode with max iterations
  MODE="build"
  PROMPT_FILE="PROMPT_build.md"
  MAX_ITERATIONS=$1
else
  # Build mode, unlimited (no arguments or invalid input)
  MODE="build"
  PROMPT_FILE="PROMPT_build.md"
  MAX_ITERATIONS=0
fi

LOG_FILE=".ralph/logs/activity_$MODE.txt"
LOG_DIR=$(dirname "$LOG_FILE")

# Make sure log directory exists
if [ ! -d "$LOG_DIR" ]; then
  mkdir -p "$LOG_DIR"
fi

ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

(
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "Mode:   $MODE"
  echo "Prompt: $PROMPT_FILE"
  echo "Branch: $CURRENT_BRANCH"
  [[ $MAX_ITERATIONS -gt 0 ]] && echo "Max:    $MAX_ITERATIONS iterations"
  echo "Date:   $(date)"
  echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
) | tee -a "$LOG_FILE"

# Verify prompt file exists
if [ ! -f "$PROMPT_FILE" ]; then
  echo "Error: $PROMPT_FILE not found"
  exit 1
fi

LAST_EC=0
while true; do
  if [[ $MAX_ITERATIONS -gt 0  && $ITERATION -ge $MAX_ITERATIONS ]]; then
    echo "Reached max iterations: $MAX_ITERATIONS"
    break
  fi

  # Run Ralph iteration with selected prompt
  # --yolo: Auto-approve all tool calls
  # --output-format=stream-json: Structured output for logging/monitoring
  # --model auto: Allow Gemini to auto-route tasks to the most appropriate model based on
  #               complexity: typically Pro for planning and Flash for building
  cat "$PROMPT_FILE" | gemini \
    --yolo \
    --output-format=stream-json \
    --model auto |
  tee -a "$LOG_FILE" |
  ./.ralph/format_output.sh
  EXIT_CODE=$PIPESTATUS[0]

  if [[ $EXIT_CODE -ne 0 && $LAST_EC -ne 0 ]]; then
    # If EXIT_CODE and LAST_EC are both non-zero then we have errored out twice
    # in a row, so it's probably time to just get out.
    echo "Agent exited with an error twice in a row"
    break
  fi
  LAST_EC=$EXIT_CODE

  # Push changes after each iteration
  git push origin "$CURRENT_BRANCH" || {
    echo "Failed to push. Creating remote branch..."
    git push -u origin "$CURRENT_BRANCH"
  }

  ITERATION=$((ITERATION + 1))
  echo -e "\n======================== ITERATION $ITERATION ========================\n" |
  tee -a "$LOG_FILE"
done

if [[ $EXIT_CODE -ne 0 ]]; then
  echo -e "\n[LOOP TERMINATED] Too many agent failures\n\n" >> "$LOG_FILE"
else
  echo -e "\n[LOOP COMPLETED]\n\n" >> "$LOG_FILE"
fi
