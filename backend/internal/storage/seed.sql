-- Seed data for development and testing

-- Demo user
INSERT OR IGNORE INTO users (id, nickname) VALUES (1, 'demo_player');

-- Problem: Two Sum
INSERT OR IGNORE INTO problems (id, title, description, time_limit_ms, memory_limit_mb) VALUES (
    1,
    'Two Sum',
    '## Problem

Given an array of integers `nums` and an integer `target`, return the indices of the two numbers that add up to `target`.

You may assume that each input has **exactly one solution**, and you may not use the same element twice.

Return the answer as two space-separated indices (0-indexed).

## Input Format

- First line: space-separated integers (the array)
- Second line: the target integer

## Output Format

- Two space-separated indices

## Example

**Input:**
```
2 7 11 15
9
```

**Output:**
```
0 1
```

**Explanation:** nums[0] + nums[1] = 2 + 7 = 9

## Constraints

- 2 ≤ nums.length ≤ 10^4
- -10^9 ≤ nums[i] ≤ 10^9
- -10^9 ≤ target ≤ 10^9
- Only one valid answer exists.',
    2000,
    256
);

-- Sample test cases (visible to users)
INSERT OR IGNORE INTO test_cases (id, problem_id, input, expected_output, is_sample) VALUES
(1, 1, '2 7 11 15
9', '0 1', 1),
(2, 1, '3 2 4
6', '1 2', 1);

-- Hidden test cases
INSERT OR IGNORE INTO test_cases (id, problem_id, input, expected_output, is_sample) VALUES
(3, 1, '3 3
6', '0 1', 0),
(4, 1, '-1 -2 -3 -4 -5
-8', '2 4', 0),
(5, 1, '1 2 3 4 5 6 7 8 9 10
19', '8 9', 0);
