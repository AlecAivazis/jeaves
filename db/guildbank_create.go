// Code generated by entc, DO NOT EDIT.

package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlecAivazis/jeeves/db/bankitem"
	"github.com/AlecAivazis/jeeves/db/guildbank"
	"github.com/facebookincubator/ent/dialect/sql"
)

// GuildBankCreate is the builder for creating a GuildBank entity.
type GuildBankCreate struct {
	config
	channelID        *string
	displayMessageID *string
	balance          *int
	items            map[int]struct{}
	guild            map[int]struct{}
}

// SetChannelID sets the channelID field.
func (gbc *GuildBankCreate) SetChannelID(s string) *GuildBankCreate {
	gbc.channelID = &s
	return gbc
}

// SetDisplayMessageID sets the displayMessageID field.
func (gbc *GuildBankCreate) SetDisplayMessageID(s string) *GuildBankCreate {
	gbc.displayMessageID = &s
	return gbc
}

// SetBalance sets the balance field.
func (gbc *GuildBankCreate) SetBalance(i int) *GuildBankCreate {
	gbc.balance = &i
	return gbc
}

// SetNillableBalance sets the balance field if the given value is not nil.
func (gbc *GuildBankCreate) SetNillableBalance(i *int) *GuildBankCreate {
	if i != nil {
		gbc.SetBalance(*i)
	}
	return gbc
}

// AddItemIDs adds the items edge to BankItem by ids.
func (gbc *GuildBankCreate) AddItemIDs(ids ...int) *GuildBankCreate {
	if gbc.items == nil {
		gbc.items = make(map[int]struct{})
	}
	for i := range ids {
		gbc.items[ids[i]] = struct{}{}
	}
	return gbc
}

// AddItems adds the items edges to BankItem.
func (gbc *GuildBankCreate) AddItems(b ...*BankItem) *GuildBankCreate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return gbc.AddItemIDs(ids...)
}

// SetGuildID sets the guild edge to Guild by id.
func (gbc *GuildBankCreate) SetGuildID(id int) *GuildBankCreate {
	if gbc.guild == nil {
		gbc.guild = make(map[int]struct{})
	}
	gbc.guild[id] = struct{}{}
	return gbc
}

// SetNillableGuildID sets the guild edge to Guild by id if the given value is not nil.
func (gbc *GuildBankCreate) SetNillableGuildID(id *int) *GuildBankCreate {
	if id != nil {
		gbc = gbc.SetGuildID(*id)
	}
	return gbc
}

// SetGuild sets the guild edge to Guild.
func (gbc *GuildBankCreate) SetGuild(g *Guild) *GuildBankCreate {
	return gbc.SetGuildID(g.ID)
}

// Save creates the GuildBank in the database.
func (gbc *GuildBankCreate) Save(ctx context.Context) (*GuildBank, error) {
	if gbc.channelID == nil {
		return nil, errors.New("db: missing required field \"channelID\"")
	}
	if gbc.displayMessageID == nil {
		return nil, errors.New("db: missing required field \"displayMessageID\"")
	}
	if gbc.balance == nil {
		v := guildbank.DefaultBalance
		gbc.balance = &v
	}
	if len(gbc.guild) > 1 {
		return nil, errors.New("db: multiple assignments on a unique edge \"guild\"")
	}
	return gbc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (gbc *GuildBankCreate) SaveX(ctx context.Context) *GuildBank {
	v, err := gbc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (gbc *GuildBankCreate) sqlSave(ctx context.Context) (*GuildBank, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(gbc.driver.Dialect())
		gb      = &GuildBank{config: gbc.config}
	)
	tx, err := gbc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(guildbank.Table).Default()
	if value := gbc.channelID; value != nil {
		insert.Set(guildbank.FieldChannelID, *value)
		gb.ChannelID = *value
	}
	if value := gbc.displayMessageID; value != nil {
		insert.Set(guildbank.FieldDisplayMessageID, *value)
		gb.DisplayMessageID = *value
	}
	if value := gbc.balance; value != nil {
		insert.Set(guildbank.FieldBalance, *value)
		gb.Balance = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(guildbank.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	gb.ID = int(id)
	if len(gbc.items) > 0 {
		p := sql.P()
		for eid := range gbc.items {
			p.Or().EQ(bankitem.FieldID, eid)
		}
		query, args := builder.Update(guildbank.ItemsTable).
			Set(guildbank.ItemsColumn, id).
			Where(sql.And(p, sql.IsNull(guildbank.ItemsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(gbc.items) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"items\" %v already connected to a different \"GuildBank\"", keys(gbc.items))})
		}
	}
	if len(gbc.guild) > 0 {
		eid := keys(gbc.guild)[0]
		query, args := builder.Update(guildbank.GuildTable).
			Set(guildbank.GuildColumn, eid).
			Where(sql.EQ(guildbank.FieldID, id).And().IsNull(guildbank.GuildColumn)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(gbc.guild) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"guild\" %v already connected to a different \"GuildBank\"", keys(gbc.guild))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return gb, nil
}
