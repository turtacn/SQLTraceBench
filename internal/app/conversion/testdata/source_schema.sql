-- a simplified version of the warehouse table from TPC-C
CREATE TABLE warehouse (
  w_id INT PRIMARY KEY,
  w_name VARCHAR(10),
  w_tax DECIMAL(4,2)
) ENGINE=OLAP
DISTRIBUTED BY HASH(w_id);
