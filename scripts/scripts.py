#!/usr/bin/python3
import MySQLdb
import MySQLdb.cursors


def get_public_attrs(o):
    return [a for a in dir(o) if not a.startswith("__")]


add_comic_format = """INSERT INTO `porn_comic_fa`.`comics` (`created_at`, `updated_at`, `name`, `description`, `status`, `number_of_episodes`, `cover_path`, `last_episode_time`, `user_id`) VALUES ('2022-07-25 11:57:46.329', '2022-07-25 11:57:46.329', 'dsfsd', 'sdfsd', 'aaa', '324', 'sdf', '2022-07-25 11:57:46.312', '1');"""

if __name__ == "__main__":
    db = MySQLdb.Connection(host="localhost", user="me",
                        passwd="@/Fucky9909", db="porn_comic_fa")

    cur = db.cursor()
    cur.execute("select * from comics")
    
    print("connected to database")

    db.close()