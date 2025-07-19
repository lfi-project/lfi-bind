int cb(int (*fn)(void)) {
    return 10 + fn();
}
