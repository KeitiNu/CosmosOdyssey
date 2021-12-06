Create table pricelist(
	id varchar(50) Primary key,
	valid_until timestamp with time zone NOT NULL,
	created timestamp with time zone NOT NULL
);

Create table route_info(
	id varchar(50) Primary key,
	pricelist_id varchar(50) NOT NULL,
	from_location varchar(50) NOT NULL,
	to_location varchar(50) NOT NULL,
	distance bigint NOT NULL,
	CONSTRAINT fk_pricelist
      FOREIGN KEY(pricelist_id) 
	  	REFERENCES pricelist(id)
	  		ON DELETE CASCADE
);

Create table reservation(
	id varchar(50) Primary key,
	pricelist_id varchar(50) NOT NULL,
	first_name varchar(50) NOT NULL,
	last_name varchar(50) NOT NULL,
	total_price decimal NOT NULL,
	total_time interval DAY TO MINUTE NOT NULL,
	flights varchar(50)[] NOT NULL,
	CONSTRAINT fk_pricelist
      FOREIGN KEY(pricelist_id) 
	  	REFERENCES pricelist(id)
	  		ON DELETE CASCADE
);

CREATE TABLE company(
	id varchar(50) Primary Key,
	name varchar(50) NOT NULL,
	pricelist_id varchar(50) NOT NULL,
	CONSTRAINT fk_pricelist
      FOREIGN KEY(pricelist_id) 
	  	REFERENCES pricelist(id)
	  		ON DELETE CASCADE
);

CREATE TABLE flight(
	id varchar(50) Primary Key,
	route_info_id varchar(50) NOT NULL,
	company_id varchar(50) NOT NULL,
	price decimal NOT NULL,
	start_time timestamp with time zone NOT NULL,
	end_time timestamp with time zone NOT NULL,

	CONSTRAINT fk_route_info
      FOREIGN KEY(route_info_id) 
	  	REFERENCES route_info(id)
	  		ON DELETE CASCADE,
	CONSTRAINT fk_company
      FOREIGN KEY(company_id) 
	  	REFERENCES company(id)
	  		ON DELETE CASCADE
);