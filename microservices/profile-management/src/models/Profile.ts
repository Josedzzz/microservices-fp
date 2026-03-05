import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn } from "typeorm";

@Entity("profiles")
export class Profile {
  @PrimaryGeneratedColumn("uuid")
  id!: string;

  @Column({ unique: true })
  employeeId!: string;

  @Column()
  name!: string;

  @Column()
  email!: string;

  @Column({ default: "" })
  phone!: string;

  @Column({ default: "" })
  address!: string;

  @Column({ default: "" })
  city!: string;

  @Column({ type: "text", default: "" })
  bio!: string;

  @CreateDateColumn()
  createdAt!: Date;
}
