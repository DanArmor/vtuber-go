-- Modify "queue_tasks" table
ALTER TABLE "queue_tasks" ADD COLUMN "worker_name" character varying NOT NULL;
