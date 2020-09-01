package database

type user struct {
	Name     string
	Email    string
	Password string
}
type technician struct {
	Name     string
	Code     int
	Email    string
	Password string
}

func initData() {
	//create function for atomatic updated_at
	Conn.Exec(`
		CREATE OR REPLACE FUNCTION trigger_set_timestamp()
		RETURNS TRIGGER AS $$
		BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)

	Conn.Exec(`
		CREATE TABLE public.users
		(
			id serial NOT NULL,
			name character varying(50) NOT NULL,
			email character varying(50) NOT NULL,
			password character varying(32) NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT now(),
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			PRIMARY KEY (id),
			CONSTRAINT users_email_uq UNIQUE (email)
		);
	`)

	Conn.Exec(`
		CREATE TABLE public.technicians
		(
			id serial NOT NULL,
			code integer NOT NULL,
			name character varying(50) NOT NULL,
			email character varying(50) NOT NULL,
			password character varying(32) NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT now(),
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			PRIMARY KEY (id)
		);
	`)

	Conn.Exec(`
		CREATE TABLE public.tickets
		(
			id serial NOT NULL,
			status character varying(8) NOT NULL,
			tracking character(32) NOT NULL,
			stars smallint NOT NULL,
			users_id integer NOT NULL,
			technicians_id integer NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT now(),
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			PRIMARY KEY (id),
			CONSTRAINT tickets_users_id FOREIGN KEY (users_id)
				REFERENCES public.users (id) MATCH SIMPLE
				ON UPDATE CASCADE
				ON DELETE RESTRICT
				NOT VALID,
			CONSTRAINT tickets_technicians_id FOREIGN KEY (technicians_id)
				REFERENCES public.technicians (id) MATCH SIMPLE
				ON UPDATE CASCADE
				ON DELETE RESTRICT
				NOT VALID
		);
	`)

	//Create the trigger
	Conn.Exec(`
		CREATE TRIGGER set_timestamp
		BEFORE UPDATE ON public.tickets
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_timestamp();
	`)

	// 	Conn.Exec(`
	// 	CREATE TABLE public.service_rating
	// 	(
	// 		id serial NOT NULL,
	// 		stars integer NOT NULL,
	// 		tickets_id integer NOT NULL,
	// 		created_at timestamp with time zone NOT NULL DEFAULT now(),
	// 		PRIMARY KEY (id),
	// 		CONSTRAINT service_rating_tickets_id FOREIGN KEY (tickets_id)
	// 			REFERENCES public.service_rating (id) MATCH SIMPLE
	// 			ON UPDATE CASCADE
	// 			ON DELETE CASCADE
	// 			NOT VALID
	// 	);
	// `)

	var users = []user{
		{Name: "Anderson", Email: "anderson@dominio.com", Password: "012"},
		{Name: "Linda", Email: "linda@dominio.com", Password: "234"},
		{Name: "Lucia", Email: "lucia@dominio.com", Password: "456"},
		{Name: "Georgi", Email: "georgi@dominio.com", Password: "678"},
	}
	Conn.Create(&users)

	var technicians = []technician{
		{Name: "David", Email: "david@harper.com", Password: "210", Code: 3323},
		{Name: "Norman", Email: "norman@harper.com", Password: "432", Code: 8112},
		{Name: "Lizeth", Email: "lizeth@harper.com", Password: "654", Code: 4211},
		{Name: "Azul", Email: "azul@harper.com", Password: "876", Code: 1017},
		{Name: "Lorena", Email: "lorena@harper.com", Password: "098", Code: 2201},
	}
	Conn.Create(&technicians)

}
