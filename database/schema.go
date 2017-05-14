package database

import "log"

// SetUpSchema creates any necessary tables if they do not already exist.
func SetUpSchema(client DBClient, logger *log.Logger) {
	rows, err := client.Connection().Query(`SELECT EXISTS(
	    SELECT *
	    FROM information_schema.tables
	    WHERE
	      table_schema = 'public' AND
	      table_name = 'jobs'
	);`)
	check(err, logger)

	var exists bool
	if rows.Next() {
		err := rows.Scan(&exists)
		check(err, logger)

		logger.Printf("Querying if jobs table exists: '%v'", exists)
	}

	if !exists {
		_, err := client.Connection().Exec(createJobsTable)
		check(err, logger)
		_, err = client.Connection().Exec(createListenFunction)
		check(err, logger)
		_, err = client.Connection().Exec(createJobsTrigger)
		check(err, logger)

		logger.Println("'jobs' table created")
	}
}

const createJobsTable string = `CREATE TABLE jobs(
    JOB_ID SERIAL PRIMARY KEY NOT NULL,
    NAME TEXT NOT NULL,
    URL TEXT NOT NULL,
	START_TIME BIGINT NOT NULL,
    INTERVAL BIGINT NOT NULL,
    CREATED_AT TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc'));`

const createListenFunction string = `CREATE OR REPLACE FUNCTION watch_jobs() RETURNS TRIGGER AS $$

	    DECLARE
	        data json;
	        notification json;

	    BEGIN

	        -- Convert the old or new row to JSON, based on the kind of action.
	        -- Action = INSERT or UPDATE?   -> NEW row
	        data = row_to_json(NEW);

	        -- Contruct the notification as a JSON string.
	        notification = json_build_object(
	                          'table',TG_TABLE_NAME,
	                          'action', TG_OP,
	                          'data', data);


	        -- Execute pg_notify(channel, notification)
	        PERFORM pg_notify('watch_tasks',notification::text);

	        -- Result is ignored since this is an AFTER trigger
	        RETURN NULL;
	    END;

	$$ LANGUAGE plpgsql;
`

const createJobsTrigger string = `CREATE TRIGGER jobs_notify_event
AFTER INSERT ON jobs
FOR EACH ROW EXECUTE PROCEDURE watch_jobs();`
