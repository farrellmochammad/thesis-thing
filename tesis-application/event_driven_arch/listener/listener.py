from pgnotify import await_pg_notifications

for notification in await_pg_notifications(
        'postgresql://postgres:postgrespassword@localhost:5432',
        ['order_progress_event']):

    print(notification.channel)
    print(notification.payload)