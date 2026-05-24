# Конфигурационные параметры (можно переопределять через -v)
BEGIN {
    # Цвета
    color_title = "\033[1;36m";
    color_project = "\033[1;37m";
    color_success = "\033[1;32m";
    color_error = "\033[1;31m";
    color_warning = "\033[1;33m";
    color_time = "\033[1;34m";
    color_eta = "\033[1;35m";
    color_reset = "\033[0m";

    # Символы
    char_done = "━";
    char_pending = "━";
    char_ok = "✅";
    char_error = "✖";
    char_warning = "⚠";
    char_time = "⏱"
    char_eta = "⌛"

    # Размеры
    bar_length = 50;
}

function begin(title, project) {
    printf "\n%s%s%s [%s%s%s]\n", color_title, title, color_reset, color_project, project, color_reset;
    start_time = systime();
}

function progress(current, total) {
    percent = int((current/total)*100);
    filled = int(bar_length*current/total);
    elapsed = systime()-start_time;

    if (current > 1) {
        remaining = int(elapsed*(total-current)/current);
        remaining_str = sprintf("%02d:%02d", remaining/60, remaining%60);
    } else {
        remaining_str = "--:--";
    }

    elapsed_str = sprintf("%02d:%02d", elapsed/60, elapsed%60);

    printf "\r[";
    for (i=1; i<=bar_length; i++) {
        if (i <= filled)
            printf "%s%s%s", color_success, char_done, color_reset;
        else
            printf "%s%s%s", color_error, char_pending, color_reset;
    }
    printf "] %s%d%%%s (%d/%d) %s%s%s%s %s%s%s%s",
        color_warning, percent, color_reset,
        current, total,
        color_time, char_time, elapsed_str, color_reset,
        color_eta, char_eta, remaining_str, color_reset;

    fflush();
}

function finish(errors, total) {
    if (errors == 0) {
        progress(total, total);
        printf "\n%s%s Done in %d seconds!%s\n", color_success, char_ok, systime()-start_time, color_reset;
    } else {
        printf "\n%s%s %d Errors detected%s", color_error, char_error, errors, color_reset;
        printf "\n%s%s Stopped after %d seconds%s\n", color_warning, char_warning, systime()-start_time, color_reset;
    }
}

# Основные обработчики
BEGIN { begin(bar_title, project) }
/^mkdir/ { progress(step++, total_steps) }
/.go:/ { errors++ }
END { finish(errors, total_steps) }
